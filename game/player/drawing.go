package player

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func DrawRotatedRect(screen *ebiten.Image, cx, cy, w, h float64, angleDeg float64, clr color.RGBA) {
	angleReg := angleDeg * math.Pi / 180

	rectImg := ebiten.NewImage(int(w), int(h))
	rectImg.Fill(clr)

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Rotate(angleReg)
	op.GeoM.Translate(cx, cy)

	screen.DrawImage(rectImg, op)
}
