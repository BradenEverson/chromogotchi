import { RequestObject, ResponseObject } from "./request";

type PetId = string;
type ErrorString = string;

var connected = false;
var id: string | null = null;

const HTTP_LOCATION: string = "http://localhost:7878";

const NAME: HTMLElement = document.getElementById("petName") as HTMLElement;
const CANVAS: HTMLCanvasElement = document.getElementById("petCanvas") as HTMLCanvasElement;
CANVAS.width = 48;
CANVAS.height = 48;

var hunger: number, sleep: number, happiness: number;
var sprite: Uint8Array;

var x = CANVAS.width / 2
var y = CANVAS.height / 2

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
        querySomething("Get")
    }
});

socket.addEventListener("message", (event) => {
    console.log(event.data);
    const data = JSON.parse(event.data) as ResponseObject;
    const dataToUint8Array = base64ToUint8Array(data.data);
    switch (data.type) {
        case "Pet":
            let name: string = uint8ArrayToString(dataToUint8Array);
            NAME.innerText = name;
            break;
        case "Sprite":
            sprite = dataToUint8Array;
            break;
        case "Slept":
            sleep = bytesToF32([
                dataToUint8Array[0],
                dataToUint8Array[1],
                dataToUint8Array[2],
                dataToUint8Array[3],
            ]);
            break;
        case "Happy":
            happiness = bytesToF32([
                dataToUint8Array[0],
                dataToUint8Array[1],
                dataToUint8Array[2],
                dataToUint8Array[3],
            ]);
            break;
        case "Fed":
            hunger = bytesToF32([
                dataToUint8Array[0],
                dataToUint8Array[1],
                dataToUint8Array[2],
                dataToUint8Array[3],
            ]);
        default:
            break;
    }
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
    let byteSizedId = encoder.encode(id);
    console.log(byteSizedId);
    let query: RequestObject = {
        type: "Establish",
        data: Array.from(byteSizedId),
    };
    sendQuery(query);
    querySomething("Play")
    querySomething("Feed")
    querySomething("Sleep")
    querySomething("Get")
    querySomething("Sprite")
}

function querySomething(thing: string) {
    let query: RequestObject = { type: thing, data: [] };
    sendQuery(query);
}

function sendQuery(query: RequestObject) {
    let queryStr = JSON.stringify(query);
    console.log(queryStr);
    socket.send(queryStr);
}
id = getCookie("PETID");
if (!id) {
    id = await newPet();
}

setInterval(() => {
    if (connected) {
        querySomething("Play")
        querySomething("Feed")
        querySomething("Sleep")
        querySomething("Sprite")

    }
}, 5000);

let ctx = CANVAS.getContext("2d")
setInterval(() => {
    if (connected && sprite && ctx) {
        /// Update the dude's position on the canvas
        ctx.clearRect(0, 0, CANVAS.width, CANVAS.height)

        const imageData = ctx.createImageData(16, 16);

        for (let i = 0; i < sprite.length; i++) {
            imageData.data[i] = sprite[i];
        }

        ctx.clearRect(0, 0, CANVAS.width, CANVAS.height);

        ctx.putImageData(imageData, x, y);
    }
}, 10);

export function base64ToUint8Array(base64: string): Uint8Array {
    const binaryString = atob(base64);
    const len = binaryString.length;
    const bytes = new Uint8Array(len);

    for (let i = 0; i < len; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }

    return bytes;
}

export function uint8ArrayToString(uint8Array: Uint8Array): string {
    const decoder = new TextDecoder("utf-8");
    return decoder.decode(uint8Array);
}

export function bytesToF32(bytes: [number, number, number, number]): number {
    let buffer = new ArrayBuffer(4);
    let view = new DataView(buffer);

    for (let index = 0; index < bytes.length; index++) {
        view.setUint8(index, bytes[index]);
    }

    return view.getFloat32(0, true);
}

export function f32ToBytes(float: number): [number, number, number, number] {
    let buffer = new ArrayBuffer(4);
    let view = new DataView(buffer);

    view.setFloat32(0, float);

    let first = view.getUint8(0);
    let second = view.getUint8(1);
    let third = view.getUint8(2);
    let fourth = view.getUint8(3);

    return [first, second, third, fourth];
}
