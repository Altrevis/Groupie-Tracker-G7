function obtenirArtistes() {
    fetch('/artist')
        .then(response => response.json())
        .then(artistes => {
            document.getElementById('listeArtistes').innerHTML = '';
            artistes.forEach(artiste => {
                const elementArtiste = document.createElement('div');
                elementArtiste.classList.add('slot');
                elementArtiste.innerHTML = `
                    <h2>${artiste.name}</h2>
                    <img src="${artiste.image}" alt="${artiste.name}">
                    <a href="artist.html?id=${artiste.id}">voir plus</a>
                `;
                document.getElementById('listeArtistes').appendChild(elementArtiste);
            });
        })
        .catch(error => console.error('Erreur lors de la récupération des artistes :', error));
}

function obtenirArtiste(artistId) {
    fetch(`/artist?id=${artistId}`)
        .then(response => response.json())
        .then(artiste => {
            document.getElementById('artist-name').textContent = artiste.name;
            document.getElementById('artist-image').src = artiste.image;
            document.getElementById('creation-date').textContent = artiste.creationDate;
            document.getElementById('first-album').textContent = artiste.firstAlbum;

            const membersList = document.getElementById('members-list');
            membersList.innerHTML = '';
            artiste.members.forEach(member => {
                const li = document.createElement('li');
                li.textContent = member;
                membersList.appendChild(li);
            });

            const concertDates = document.getElementById('concert-dates');
            concertDates.innerHTML = '';
            for (const date in artiste.concertDates) {
                const li = document.createElement('li');
                li.classList.add('collection-item');
                li.innerHTML = `
                    ${date}
                    <span class="secondary-content location">${artiste.concertDates[date]}</span>
                `;
                concertDates.appendChild(li);
            }
        })
        .catch(error => console.error('Erreur lors de la récupération des détails de l\'artiste :', error));
}

window.onload = function() {
    const urlParams = new URLSearchParams(window.location.search);
    const artistId = urlParams.get('id');

    if (artistId) {
        obtenirArtiste(artistId);
    } else {
        obtenirArtistes();
    }
};