const fs = require("fs");

const PATH = "../dataset";
const TRACK_LIST = "../output/track-list.json";
const DATAS_OUTPUT = "../output/playlists_.json";

const tracks = JSON.parse(fs.readFileSync(TRACK_LIST));
let playlists = [];

let countSum = 0;

function openFile(path) {
    const datas = JSON.parse(fs.readFileSync(path));
    datas["playlists"].map((playlist) => {

        const myPlaylist = [];

        playlist["tracks"].map((track) => {
            const uri = track["track_uri"];

            if (typeof tracks[uri] != "undefined") {
                myPlaylist.push(tracks[uri].id);
            }
        });
        

        // every playlist should contains at least 3 tracks from the database
        if (myPlaylist.length > 2) {
            playlists.push(myPlaylist);
            countSum += myPlaylist.length;
        }
    });
    console.log(path, "total playlist number:", playlists.length)
}

const files = fs.readdirSync(PATH);
for (let i=0; i<files.length; i++) {
    openFile(PATH + "/" + files[i])
}

console.log("Average track number by playlist: ", countSum / playlists.length)
console.log("Playlist count", playlists.length)

fs.writeFileSync(DATAS_OUTPUT, JSON.stringify(playlists))