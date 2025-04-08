package main

// Serie representa la estructura de datos para una serie de TV en la aplicación.
// Contiene los campos que se mapean a las columnas de la tabla 'series'
// y las etiquetas JSON para la serialización/deserialización.
type Serie struct {
	ID                 int    `json:"id"`
	Title              string `json:"title"`
	Status             string `json:"status"`
	LastEpisodeWatched int    `json:"lastEpisodeWatched"`
	TotalEpisodes      int    `json:"totalEpisodes"`
	Ranking            int    `json:"ranking"`
}
