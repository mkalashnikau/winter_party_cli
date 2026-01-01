package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

var (
	seasonalRunes = []rune{'❄', '❅', '❆', '✨', '🎄'}
	gentleRunes   = seasonalRunes[:3]
)

type particle struct {
	x  float64
	y  float64
	vx float64
	vy float64
	ch rune
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("unable to create screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("unable to init screen: %v", err)
	}
	defer screen.Fini()

	rand.Seed(time.Now().UnixNano())

	frameDuration := time.Second / 60
	ticker := time.NewTicker(frameDuration)
	defer ticker.Stop()

	eventCh := make(chan tcell.Event, 10)
	// Poll events asynchronously so the animation loop can stay on schedule.
	go func() {
		for {
			event := screen.PollEvent()
			if event == nil {
				close(eventCh)
				return
			}
			eventCh <- event
		}
	}()

	flakes := make([]particle, 0, 256)
	lastSpawn := time.Now()

	runeSet := seasonalRunes

	for {
		select {
		case <-ticker.C:
			width, height := screen.Size()

			if time.Since(lastSpawn) >= 120*time.Millisecond {
				flakes = append(flakes, newFlake(width))
				lastSpawn = time.Now()
			}

			screen.Clear()

			drawHeader(screen, width)

			next := flakes[:0]
			for i := range flakes {
				flake := &flakes[i]
				flake.x += flake.vx
				flake.y += flake.vy

				if flake.y >= float64(height) {
					continue
				}
				if flake.x < 0 {
					flake.x += float64(width)
				}
				if flake.x >= float64(width) {
					flake.x -= float64(width)
				}

				x := int(flake.x)
				y := int(flake.y)
				if x >= 0 && x < width && y >= 0 && y < height {
					screen.SetContent(x, y, flake.ch, nil, tcell.StyleDefault)
					next = append(next, *flake)
				}
			}
			flakes = next

			screen.Show()
		case ev, ok := <-eventCh:
			if !ok {
				return
			}
			switch tev := ev.(type) {
			case *tcell.EventResize:
				screen.Sync()
			case *tcell.EventKey:
				switch tev.Key() {
				case tcell.KeyRune:
					switch tev.Rune() {
					case 'q', 'Q':
						return
					}
				case tcell.KeyEsc:
					spawnBurst(&flakes, screen, runeSet)
				case tcell.KeyCtrlC:
					return
				}
			}
		}
	}
}

func newFlake(width int) particle {
	x := rand.Float64() * float64(width)
	vy := 0.25 + rand.Float64()*0.6
	vx := (rand.Float64() - 0.5) * 0.2
	runes := gentleRunes
	return particle{
		x:  x,
		y:  0,
		vx: vx,
		vy: vy,
		ch: runes[rand.Intn(len(runes))],
	}
}

func spawnBurst(flakes *[]particle, screen tcell.Screen, runeSet []rune) {
	width, _ := screen.Size()
	for i := 0; i < 32; i++ {
		x := rand.Float64() * float64(width)
		vy := 0.35 + rand.Float64()*0.9
		vx := (rand.Float64() - 0.5) * 0.4
		ch := runeSet[rand.Intn(len(runeSet))]
		*flakes = append(*flakes, particle{
			x:  x,
			y:  0,
			vx: vx,
			vy: vy,
			ch: ch,
		})
	}
}

func drawHeader(screen tcell.Screen, width int) {
	header := "Happy Holidays! Press Esc for more snow, q to quit."
	x := (width - len(header)) / 2
	if x < 0 {
		x = 0
	}
	for i, r := range header {
		screen.SetContent(x+i, 0, r, nil, tcell.StyleDefault)
	}
}
