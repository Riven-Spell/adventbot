package utils

import (
	"encoding/json"
	"os"
)

func SaveFile(name string, data interface{}) {
	if _, err := os.Stat(name); err == nil {
		f, _ := os.OpenFile(name, os.O_RDWR|os.O_TRUNC, 0644)
		d, _ := json.Marshal(data)
		_,_ = f.Write(d) //Yes I know these have unhandled errors.
		_ = f.Close()
		println("Opened file and wrote", name)
		return
	} else {
		if f, err := os.Create(name); err == nil {
			d, _ := json.Marshal(data)
			_, _ = f.Write(d) //Yes I know these have unhandled errors.
			_ = f.Close()
			println("Wrote file", name)
			return
		} else {
			println("Couldn't create file", name)
			return
		}
	}
}
