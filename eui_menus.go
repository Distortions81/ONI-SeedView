package main

import (
	eui "EUI/eui"
)

func (g *Game) initEUI() {
	eui.SetUIScale(float32(uiScale))

	w := eui.NewWindow(&eui.WindowData{
		Title:     ScreenshotMenuTitle,
		Size:      eui.Point{X: 160, Y: 120},
		Open:      false,
		Resizable: false,
	})
	flow := &eui.ItemData{ItemType: eui.ITEM_FLOW, Size: w.Size, FlowType: eui.FLOW_VERTICAL}
	w.AddItem(flow)

	for i, q := range ScreenshotQualities {
		btn, ev := eui.NewButton(&eui.ItemData{Text: q, Size: eui.Point{X: 140, Y: 20}, FontSize: 8})
		idx := i
		ev.Handle = func(ev eui.UIEvent) {
			if ev.Type == eui.EventClick {
				g.ssQuality = idx
			}
		}
		flow.AddItem(btn)
	}

	bwBtn, bwEv := eui.NewButton(&eui.ItemData{Text: ScreenshotBWLabel, Size: eui.Point{X: 140, Y: 20}, FontSize: 8})
	bwEv.Handle = func(ev eui.UIEvent) {
		if ev.Type == eui.EventClick {
			g.ssNoColor = !g.ssNoColor
			g.noColor = g.ssNoColor
		}
	}
	flow.AddItem(bwBtn)

	saveBtn, saveEv := eui.NewButton(&eui.ItemData{Text: ScreenshotSaveLabel, Size: eui.Point{X: 140, Y: 20}, FontSize: 8})
	saveEv.Handle = func(ev eui.UIEvent) {
		if ev.Type == eui.EventClick {
			if g.ssPending == 0 {
				g.ssPending = 2
			}
		}
	}
	flow.AddItem(saveBtn)

	cancelBtn, cancelEv := eui.NewButton(&eui.ItemData{Text: ScreenshotCancelLabel, Size: eui.Point{X: 140, Y: 20}, FontSize: 8})
	cancelEv.Handle = func(ev eui.UIEvent) {
		if ev.Type == eui.EventClick {
			g.showShotMenu = false
			w.Open = false
		}
	}
	flow.AddItem(cancelBtn)

	g.ssWindow = w
	w.AddWindow(false)
}
