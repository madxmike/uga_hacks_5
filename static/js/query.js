let map = L.map('mapid').setView([51.505, -0.09], 13);

L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
}).addTo(map);


var marker = L.marker([51.5, -0.09]).addTo(map);
search.onsubmit = async (e) => {
    e.preventDefault();
    let form = new FormData(search);
    form.append("Lat", map.getCenter().lat);
    form.append("Long", map.getCenter().lng);
    let response = await fetch(window.location.origin + '/api/search', {
      method: 'POST',
      body: form
    });

    let result = await response.json();

    alert(result.message);
  };