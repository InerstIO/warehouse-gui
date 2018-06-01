package main

import (
	"fmt"
	"flag"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"warehouse-optimizer/warehouse"
	"io/ioutil"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

const (
	gridPath = "../../warehouse-grid.csv"
	dimPath = "../../item-dimensions-tabbed.txt"
	batchOutput = "../../batchoutput.json"
	timeLimit = 10.0
)

// Vars
var (
	AppName string
	BuiltAt string
	debug   = flag.Bool("d", false, "enables the debug mode")
	w       *astilectron.Window
	dim	map[int][]float64 = warehouse.ParesDimensionInfo(dimPath)
	pm map[int]warehouse.Product = warehouse.ParseProductInfo(gridPath, dim)
	pathInfo map[warehouse.Point]map[warehouse.Point]float64 = warehouse.BuildPathInfo(gridPath)
	start, end warehouse.Point
	ros []warehouse.RouteOrder
)

type ose struct {
	Order warehouse.Order
	Start, End warehouse.Point
}

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "input":
		// Unmarshal payload
		var s []string
		if err = json.Unmarshal(m.Payload, &s); err != nil {
			payload = err.Error()
		    return
		} else {
			payload = findPath(s)
		}
	case "batchresult":
		var file *os.File
		file, err = os.Open(batchOutput) // For read access.
		if err != nil {
			payload = err.Error()
			return
		}
		defer file.Close()
		r, _ := ioutil.ReadAll(file)
		if err = json.Unmarshal(r, &ros); err != nil {
			payload = err.Error()
		    return
		}
		payload = string(r)
	case "route":
		// Unmarshal payload
		var s ose
		if err = json.Unmarshal(m.Payload, &s); err != nil {
			payload = err.Error()
		    return
		} else {
			var routeB []byte
			if routeB, err = json.Marshal(warehouse.Route2String(s.Order, s.Start, s.End, pm)); err != nil {
				payload = err.Error()
				return
			} else {
				payload = string(routeB)
			}
		}
	case "effort":
		// Unmarshal payload
		var s ose
		if err = json.Unmarshal(m.Payload, &s); err != nil {
			payload = err.Error()
		    return
		} else {
			var retstr string
			if effort, missWeightData := warehouse.RouteEffort(s.Order, s.Start, s.End, pm, pathInfo); missWeightData {
				retstr = fmt.Sprintf("There are some item(s) with no weight data, and the effort of this path is at least %v.\n", effort)
			} else {
				retstr = fmt.Sprintf("The effort is %v.\n", effort)
			}
			var effortB []byte
			if effortB, err = json.Marshal(retstr); err != nil {
				payload = err.Error()
				return
			} else {
				payload = string(effortB)
			}
		}
	}
	return
}

func findPath(s []string) string {
	x,_:=strconv.Atoi(s[0])
	y,_:=strconv.Atoi(s[1])
	start = warehouse.Point{X: x, Y: y}
	x,_=strconv.Atoi(s[2])
	y,_=strconv.Atoi(s[3])
	end = warehouse.Point{X: x, Y: y}
	o := strings.Split(s[5], " ")
	order := make(warehouse.Order, len(o))
	for i := range o {
		o[i] = strings.TrimSpace(o[i])
		order[i].ProdID, _ = strconv.Atoi(o[i])
		_, ok := pm[order[i].ProdID]
		if !ok {
			astilog.Fatalf("Item id %v not exist.", order[i])
		}
	}
	var optimalOrders []warehouse.Order
	orders := []warehouse.Order{order}
	if weight, err:=strconv.ParseFloat(s[4], 64); err == nil {
		orders = warehouse.SplitOrder(order, pm, weight)
	}
	for _, so := range orders {
		optimalOrders = append(optimalOrders, warehouse.BnBOrderOptimizer(so, start, end, pm, pathInfo, timeLimit))
	}
	ro := warehouse.Orders2Routes(optimalOrders, start, end, pm)
	roB, _:=json.Marshal(ro)
	return string(roB)
}

func main() {
	// Init
	flag.Parse()
	astilog.FlagInit()

	// Run bootstrap
	astilog.Debugf("Running app built at %s", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset: Asset,
		RestoreAssets:  RestoreAssets,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
		},
		Debug:    *debug,
		Homepage: "index.html",
		MenuOptions: []*astilectron.MenuItemOptions{{
			Label: astilectron.PtrStr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{Label: astilectron.PtrStr("About")},
				{Role: astilectron.MenuItemRoleClose},
			},
		}},
		OnWait: func(_ *astilectron.Astilectron, iw *astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = iw
			w.OpenDevTools()
			// This will send a message and execute a callback
			bootstrap.SendMessage(w, "route", "hello", func(m *bootstrap.MessageIn) {
				// Unmarshal payload
				var s string
				json.Unmarshal(m.Payload, &s)
				// Process message
				astilog.Infof("received %s", s)
			})
			return nil
		},
		MessageHandler: handleMessages, 
		WindowOptions: &astilectron.WindowOptions{
			BackgroundColor: astilectron.PtrStr("#333"),
			Center:          astilectron.PtrBool(true),
			Height:          astilectron.PtrInt(1000),
			Width:           astilectron.PtrInt(1600),
			WebPreferences:	&astilectron.WebPreferences{DevTools: astilectron.PtrBool(true)},
		},
	}); err == nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}