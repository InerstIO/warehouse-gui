package main

import (
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
)

// Vars
var (
	AppName string
	BuiltAt string
	debug   = flag.Bool("d", false, "enables the debug mode")
	w       *astilectron.Window
	dim	map[int][]float64 = warehouse.ParesDimensionInfo(dimPath)
	m map[int]warehouse.Product = warehouse.ParseProductInfo(gridPath, dim)
	pathInfo map[warehouse.Point]map[warehouse.Point]float64 = warehouse.BuildPathInfo(gridPath)
)

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
		payload = string(r)
	}
	return
}

func findPath(s []string) string {
	x,_:=strconv.Atoi(s[0])
	y,_:=strconv.Atoi(s[1])
	start := warehouse.Point{X: x, Y: y}
	x,_=strconv.Atoi(s[2])
	y,_=strconv.Atoi(s[3])
	end := warehouse.Point{X: x, Y: y}
	o := strings.Split(s[5], " ")
	order := make([]int, len(o))
	for i := range o {
		o[i] = strings.TrimSpace(o[i])
		order[i], _ = strconv.Atoi(o[i])
		_, ok := m[order[i]]
		if !ok {
			astilog.Fatalf("Item id %v not exist.", order[i])
		}
	}
	var optimalOrder warehouse.Order
	s4,_:=strconv.Atoi(s[4])
	if s4 == 0 {
		optimalOrder = warehouse.NNIOrderOptimizer(order, start, end, m, pathInfo)
	} else {
		optimalOrder = warehouse.BnBOrderOptimizer(order, start, end, m, pathInfo)
	}
	return string(warehouse.Routes2JSON([]warehouse.Order{optimalOrder}, start, end, m))
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
			Height:          astilectron.PtrInt(930),
			Width:           astilectron.PtrInt(810),
			WebPreferences:	&astilectron.WebPreferences{DevTools: astilectron.PtrBool(true)},
		},
	}); err == nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}