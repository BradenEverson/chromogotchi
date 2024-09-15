import { RequestObject } from "./request";

type PetId = string;
type ErrorString = string;

var connected = false;
var id: string | null = null;

const HTTP_LOCATION: string = "http://localhost:7878";

export function getCookie(name: string): PetId | null {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop()?.split(";").shift() || null;
    return null;
}

export async function newPet(): Promise<PetId | ErrorString> {
    try {
        const response = await fetch(HTTP_LOCATION + "/newpet", {
            method: "GET",
        });

        if (response.status === 200) {
            const data = await response.json();
            const petID = data.PetId;
            document.cookie = `PETID=${petID}; path=/;`;
            return petID as PetId;
        } else {
            return "Bad response from endpoint :(" as ErrorString;
        }
    } catch (err) {
        return "Error reaching HTTP endpoint" as ErrorString;
    }
}

export const socket = new WebSocket(HTTP_LOCATION + "/connection");

socket.addEventListener("open", () => {
    console.log("Connected to the WebSocket server.");
    connected = true;

    if (id) {
        console.log(id);
        establish(id);
        getPet();
    }
});

socket.addEventListener("message", (event) => {
    console.log(event.data);
    const data = JSON.parse(event.data);
    console.log(data);
});

socket.addEventListener("close", () => {
    console.log("Disconnected from the WebSocket server.");
    connected = false;
});

socket.addEventListener("error", (error) => {
    console.error("WebSocket error:", error);
});

window.addEventListener("beforeunload", () => {
    socket.close();
});

function establish(id: string) {
    let encoder = new TextEncoder();
    let byteSizedId = encoder.encode(id)
    console.log(byteSizedId)
    let query: RequestObject = { type: "Establish", data: Array.from(byteSizedId) };
    sendQuery(query);
}

function queryHunger() {
    let query: RequestObject = { type: "Feed", data: [] };
    sendQuery(query);
}

function queryHappiness() {
    let query: RequestObject = { type: "Sleep", data: [] };
    sendQuery(query);
}

function querySleep() {
    let query: RequestObject = { type: "Play", data: [] };
    sendQuery(query);
}

function getPet() {
    let query: RequestObject = { type: "Get", data: [] };
    sendQuery(query);
}

function sendQuery(query: RequestObject) {
    let queryStr = JSON.stringify(query)
    console.log(queryStr)
    socket.send(queryStr)
}
id = getCookie("PETID");
if (!id) {
    id = await newPet();
}

let hunger: number, happiness: number, sleepyness: number;
let sprite: number[];
let lilGuyX: number,
    lilGuyY: number = 0;

setInterval(() => {
    if (connected) {
        /// Get Hunger
        /// Get Happiness
        /// Get Sleepyness
        /// Get Necessary Sprite
    }

}, 1000);

setInterval(() => {
    if (connected) {
        /// Update the dude's position on the canvas
    }
}, 250);
