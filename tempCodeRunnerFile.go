func (game *Game) DrawBrushGhost(screen *ebiten.Image) {
    // Calculate the size of the brush
    radius := float64(game.BrushSize) / 2.0

    // Create a new image with size equal to the diameter of the brush plus an extra pixel.
    brushImage := ebiten.NewImage(game.BrushSize+1, game.BrushSize+1)

    // Calculate the pre-multiplied alpha for RGB
    alpha := 100 
    factor := float64(alpha) / 255
    red, green, blue := uint8(factor*255), uint8(factor*255), uint8(factor*255)

    // Iterate over the pixels of the image and color the pixels that fall inside the brush's circle.
    for row := -radius; row <= radius; row++ {
        for col := -radius; col <= radius; col++ {
            dist := math.Hypot(float64(row), float64(col))
            if dist <= radius {
                ix := int(math.Round(radius + col))
                iy := int(math.Round(radius + row))
                brushImage.Set(ix, iy, color.RGBA{red, green, blue, uint8(alpha)})
            }
        }
    }

    // Get the mouse position in the screen coordinates
    mouseX, mouseY := ebiten.CursorPosition()

    // Calculate the offset based on the radius of the brush.
    offsetX := radius * float64(PixelSize)
    offsetY := radius * float64(PixelSize)

    // Draw the brush image at the mouse position, offset by the brush radius.
    // And scale the image to the screen pixels size.
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Scale(float64(PixelSize), float64(PixelSize))
    op.GeoM.Translate(float64(mouseX)-offsetX, float64(mouseY)-offsetY)
    screen.DrawImage(brushImage, op)
}