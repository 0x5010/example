package main

import (
	"time"

	pb "gopkg.in/cheggaaa/pb.v2"
)

func main() {
	simple()
	fromPreset()
	customTemplate(`Custom template: {{counters . }}`)
	customTemplate(`{{ red "With colors:" }} {{bar . | green}} {{speed . | blue }}`)
	customTemplate(`{{ red "With funcs:" }} {{ bar . "<" "-" (cycle . "↖" "↗" "↘" "↙" ) "." ">"}} {{speed . | rndcolor }}`)
	customTemplate(`{{ bar . "[<" "·····•·····" (rnd "ᗧ" "◔" "◕" "◷" ) "•" ">]"}}`)
}

func simple() {
	count := 1000
	bar := pb.StartNew(count)
	for i := 0; i < count; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond * 2)
	}
	bar.Finish()
}

func fromPreset() {
	count := 1000
	bar := pb.Full.Start(count)
	defer bar.Finish()
	bar.Set("prefix", "fromPreset(): ")
	for i := 0; i < count/2; i++ {
		bar.Add(2)
		time.Sleep(time.Millisecond * 4)
	}
}

func customTemplate(tmpl string) {
	count := 1000
	bar := pb.ProgressBarTemplate(tmpl).Start(count)
	defer bar.Finish()
	for i := 0; i < count; i++ {
		bar.Add(1)
		time.Sleep(time.Millisecond * 16)
	}
}
