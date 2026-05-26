// reset select'a city (wyłączony)
function resetCitySelect() {
    const citySelect = document.getElementById("city");
    const opt = document.createElement("option");
    opt.value = "";
    opt.textContent = "Wybierz miasto";
    citySelect.appendChild(opt);
    citySelect.disabled = true;
}

// listener zmian w select'cie country, który pobiera miasta dla select'a city po wyborze kraju
document.getElementById("country").addEventListener("change", async function () {
    const citySelect = document.getElementById("city");
    citySelect.innerHTML = "";
    if (this.value == "") {
        resetCitySelect();
        return;
    }

    const res = await fetch("/cities?country=" + this.value);
    if (!res.ok) {
        console.log("Error:", res.status);

        resetCitySelect();
        return;
    }

    const cities = await res.json();

    // dodanie opcji w select'cie
    cities.forEach(city => {
        const opt = document.createElement("option");
        opt.value = city.Name;
        opt.textContent = city.Name;
        citySelect.appendChild(opt);
    });

    if (cities.length > 0) {
        citySelect.disabled = false;
    }
});


// tablica interpretacji kodu pogody
const weatherRanges = [
    { min: 0, max: 0, icon: "wi-day-sunny", description: "Bezchmurne niebo" },
    { min: 1, max: 1, icon: "wi-day-cloudy", description: "Lekkie zachmurzenie" },
    { min: 2, max: 2, icon: "wi-cloud", description: "Częściowe zachmurzenie" },
    { min: 3, max: 3, icon: "wi-cloudy", description: "Pełne zachmurzenie" },
    { min: 4, max: 9, icon: "wi-dust", description: "Pył lub ograniczenie widzialności" },
    { min: 10, max: 12, icon: "wi-windy", description: "Mgła przyziemna" },
    { min: 13, max: 17, icon: "wi-cloudy-gusts", description: "Zbliżające się opady lub nawałnice" },
    { min: 18, max: 19, icon: "wi-tornado", description: "Silny wiatr, możliwe tornada" },
    { min: 20, max: 29, icon: "wi-night-cloudy-windy", description: "Opady lub mgła w ciągu ostatniej godziny" },
    { min: 30, max: 39, icon: "wi-sandstorm", description: "Burza pyłowa/piaskowa lub zamieć" },
    { min: 40, max: 49, icon: "wi-fog", description: "Mgła" },
    { min: 50, max: 59, icon: "wi-sleet", description: "Mżawka" },
    { min: 60, max: 62, icon: "wi-showers", description: "Deszcz" },
    { min: 63, max: 65, icon: "wi-rain", description: "Intensywny deszcz" },
    { min: 66, max: 69, icon: "wi-snow", description: "Marznący deszcz lub deszcz zmeiszany ze śniegiem" },
    { min: 70, max: 72, icon: "wi-snow", description: "Lekki opad śniegu" },
    { min: 73, max: 75, icon: "wi-snow", description: "Intensywny opad śniegu" },
    { min: 76, max: 79, icon: "wi-snow", description: "Lodowy pył" },
    { min: 80, max: 88, icon: "wi-sleet", description: "Przelotne opady" },
    { min: 89, max: 90, icon: "wi-hail", description: "Opady gradu" },
    { min: 91, max: 94, icon: "wi-night-lightning", description: "Burza w ciągu ostatniej godziny" },
    { min: 95, max: 96, icon: "wi-storm-showers", description: "Chmury burzowe" },
    { min: 97, max: 98, icon: "wi-thunderstorm", description: "Silna burza" },
    { min: 99, max: 99, icon: "wi-hail", description: "Burza z gradem" }
];

// funkcja pobierająca informacje pogodowe dla danego miasta
async function getWeather() {
    const citySelect = document.getElementById("city");
    const cityName = citySelect.value;
    if (cityName != "") {
        const res = await fetch("/weather?city=" + cityName);
        if (!res.ok) {
            console.log("Error:", res.status);
            return;
        }

        const data = await res.json();

        // interpretacja wyników
        const dateAndTime = data.LocalTime.split("T");
        const weather = weatherRanges.find(
            r => data.WeatherCode >= r.min && data.WeatherCode <= r.max
        );

        // wyświetlenie tabeli z informacjami pogodowymi
        document.getElementById("result").innerHTML = `
            <table class="table table-striped">
                <thead>
                    <tr>
                        <th>${cityName.toUpperCase()}</th>
                        <th>${dateAndTime[0]} ${dateAndTime[1]}</th>
                    </tr>
                </thead>

                <tbody>
                    <tr>
                        <td>Pogoda</td>
                        <td><i class="wi ${weather.icon}"></i></td>
                    </tr>

                    <tr>
                        <td>Opis</td>
                        <td>${weather.description}</td>
                    </tr>

                    <tr>
                        <td>Temperatura</td>
                        <td>${data.Temperature} °C</td>
                    </tr>

                    <tr>
                        <td>Prędkość wiatru</td>
                        <td>${data.WindSpeed} km/h</td>
                    </tr>

                    <tr>
                        <td>Kierunek wiatru</td>
                        <td><i class="wi wi-wind towards-${data.WindDirection}-deg"></i></td>
                    </tr>

                    <tr>
                        <td>Pora dnia</td>
                        <td>${data.Day === 1 ? `<i class="wi wi-day-sunny"></i>` : `<i class="wi wi-night-clear"></i>`}</td>
                    </tr>

                </tbody>
            </table>
        `;
    }
}