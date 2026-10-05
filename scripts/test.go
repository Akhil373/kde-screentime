package main

func main() {
	err := load_data("/home/axle/.local/share/screentime/screen_time.db", 10)
	if err != nil {
		panic(err)
	}

	// pi := PieExamples{}
	// pi.Examples()
}
