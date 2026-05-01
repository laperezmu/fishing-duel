package encounter

type CastController struct {
	bands        []CastBand
	slotWidth    int
	totalSlots   int
	positions    []int
	positionStep int
}

func NewCastController(slotWidth int, positions []int) CastController {
	bands := OrderedCastBands()
	totalSlots := len(bands) * slotWidth
	if slotWidth <= 0 {
		slotWidth = 1
		totalSlots = len(bands)
	}
	if len(positions) == 0 {
		positions = buildCastPositions(totalSlots)
	}

	return CastController{
		bands:      bands,
		slotWidth:  slotWidth,
		totalSlots: totalSlots,
		positions:  positions,
	}
}

func (controller CastController) TotalSlots() int {
	return controller.totalSlots
}

func (controller *CastController) CurrentPosition() int {
	if len(controller.positions) == 0 {
		return 0
	}

	index := controller.positionStep
	if index >= len(controller.positions) {
		index = len(controller.positions) - 1
	}

	return controller.positions[index]
}

func (controller *CastController) Advance() int {
	if len(controller.positions) == 0 {
		return 0
	}

	controller.positionStep++
	if controller.positionStep >= len(controller.positions) {
		controller.positionStep = 0
	}

	return controller.CurrentPosition()
}

func (controller CastController) ResolveBand(position int) CastBand {
	if len(controller.bands) == 0 {
		return ""
	}
	if position < 0 {
		return controller.bands[0]
	}

	bandIndex := position / controller.slotWidth
	if bandIndex >= len(controller.bands) {
		bandIndex = len(controller.bands) - 1
	}

	return controller.bands[bandIndex]
}

func buildCastPositions(totalSlots int) []int {
	if totalSlots <= 1 {
		return []int{0}
	}

	positions := make([]int, 0, totalSlots*2-2)
	for position := 0; position < totalSlots; position++ {
		positions = append(positions, position)
	}
	for position := totalSlots - 2; position > 0; position-- {
		positions = append(positions, position)
	}

	return positions
}
