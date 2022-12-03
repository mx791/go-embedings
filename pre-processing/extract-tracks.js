const fs = require("fs");
const tracks = {};
let nextId = 0;

const PATH = "../dataset";
const OUTPUT = "../output/track-list.json";
const MAX_TRACK_COUNT = 150000;

function openFile(path) {
    const datas = JSON.parse(fs.readFileSync(path));
    let countAtStart = Object.keys(tracks).length;
    datas["playlists"].map((playlist) => {
        playlist["tracks"].map((track) => {
            const uri = track["track_uri"];

            if (typeof tracks[uri] == "undefined") {
                tracks[uri] = {
                    uri: uri,
                    id: nextId,
                    name: track['track_name'],
                    artist_name: track['artist_name'],
                    count: 1
                };
                nextId += 1;
            } else {
                tracks[uri].count += 1;
            }
        });
    });

    console.log(path, datas["playlists"].length, "playlists", Object.keys(tracks).length-countAtStart, "tracks")
}

// create a list of track, ordered by their frequencies
function sortBypopularity() {
    const tracklist = Object.keys(tracks).sort((a, b) => {
        return tracks[b].count - tracks[a].count
    });
    return tracklist;
}

// process all files from the dataset directory
const files = fs.readdirSync(PATH);
for (let i=0; i<files.length; i++) {
    openFile(PATH + "/" + files[i])
}
console.log(Object.keys(tracks).length, "tracks in the list")

// display top 10 tracks
const topTracks = sortBypopularity();
for (let i=0; i<10; i++) {
    console.log(tracks[topTracks[i]])
}
const final_tracks = {};

for (let i=0; i<MAX_TRACK_COUNT; i++) {
    const id = tracks[topTracks[i]].uri;
    final_tracks[id] = tracks[topTracks[i]];
    final_tracks[id].id = i;
}

// output in a JSON file
fs.writeFileSync(OUTPUT, JSON.stringify(final_tracks))