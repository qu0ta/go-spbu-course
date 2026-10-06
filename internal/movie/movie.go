package movie

type Movie struct {
	ID       int    `json:"id"`
	Year     int    `json:"year"`
	Title    string `json:"title"`
	Director string `json:"director"`
}
