package main

import (
	eui "EUI/eui"
	"time"
)

func (g *Game) initTray() {
	size := float32(g.iconSize())
	flow := &eui.ItemData{ItemType: eui.ITEM_FLOW, FlowType: eui.FLOW_HORIZONTAL,
		PinTo: eui.PIN_BOTTOM_RIGHT, Margin: 4,
		Size: eui.Point{X: size*4 + 16, Y: size}}

	var ssEv *eui.EventHandler
	g.ssBtn, ssEv = eui.NewButton(&eui.ItemData{Size: eui.Point{X: size, Y: size}})
	if img, ok := g.icons["../icons/camera.png"]; ok {
		g.ssBtn.Image = img
	}
	ssEv.Handle = func(ev eui.UIEvent) {
		if ev.Type == eui.EventClick {
			if time.Since(g.lastShotClick) < MenuToggleDelay {
				return
			}
			if g.showShotMenu {
				g.showShotMenu = false
				g.noColor = false
			} else {
				g.closeMenus()
				g.showShotMenu = true
				g.noColor = g.ssNoColor
			}
			g.lastShotClick = time.Now()
			g.needsRedraw = true
		}
	}
	flow.AddItem(g.ssBtn)

	var gyEv *eui.EventHandler
	g.geyserBtn, gyEv = eui.NewButton(&eui.ItemData{Size: eui.Point{X: size, Y: size}})
	if img, ok := g.icons["geyser_water.png"]; ok {
		g.geyserBtn.Image = img
	}
	gyEv.Handle = func(ev eui.UIEvent) {
		if ev.Type == eui.EventClick {
			if g.showGeyserList {
				g.showGeyserList = false
			} else {
				g.closeMenus()
				g.dragging = false
				g.showGeyserList = true
			}
			g.lastGeyserClick = time.Now()
			g.needsRedraw = true
		}
	}
	flow.AddItem(g.geyserBtn)

	var helpEv *eui.EventHandler
	g.helpBtn, helpEv = eui.NewButton(&eui.ItemData{Size: eui.Point{X: size, Y: size}})
	if img, ok := g.icons["../icons/help.png"]; ok {
		g.helpBtn.Image = img
	}
	helpEv.Handle = func(ev eui.UIEvent) {
		if ev.Type == eui.EventClick {
			if time.Since(g.lastHelpClick) < MenuToggleDelay {
				return
			}
			g.showHelp = !g.showHelp
			g.lastHelpClick = time.Now()
			g.needsRedraw = true
		}
	}
	flow.AddItem(g.helpBtn)

	var optEv *eui.EventHandler
	g.optsBtn, optEv = eui.NewButton(&eui.ItemData{Size: eui.Point{X: size, Y: size}})
	if img, ok := g.icons["../icons/gear.png"]; ok {
		g.optsBtn.Image = img
	}
	optEv.Handle = func(ev eui.UIEvent) {
		if ev.Type == eui.EventClick {
			if time.Since(g.lastOptionsClick) < MenuToggleDelay {
				return
			}
			if g.showOptions {
				g.showOptions = false
			} else {
				g.closeMenus()
				g.showOptions = true
			}
			g.lastOptionsClick = time.Now()
			g.needsRedraw = true
		}
	}
	flow.AddItem(g.optsBtn)

	g.tray = flow
	eui.AddOverlayFlow(flow)
}
