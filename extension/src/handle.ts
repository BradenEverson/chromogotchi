import { Objective, RequestObject, ResponseObject } from "./request";

type PetId = string;
type ErrorString = string;

var connected = false;
var id: string | null = null;

var goTo = [randomCoords(), randomCoords()];
var currObjective: Objective = "wander"

//const HTTP_LOCATION: string = "http://localhost:7878";
const HTTP_LOCATION: string = "https://chromogotchi-32d490e078e9.herokuapp.com";

const NAME: HTMLElement = document.getElementById("petName") as HTMLElement;
const CANVAS: HTMLCanvasElement = document.getElementById("petCanvas") as HTMLCanvasElement;
CANVAS.width = 48;
CANVAS.height = 48;

const FEED: HTMLElement = document.getElementById("feedButton") as HTMLElement;
const PLAY: HTMLElement = document.getElementById("playButton") as HTMLElement;
const SLEEP: HTMLElement = document.getElementById("sleepButton") as HTMLElement;

var sprite: Uint8Array;

var x = randomCoords();
var y = randomCoords();

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
        querySomething("Sprite")
    }
}, 15000);

let ctx = CANVAS.getContext("2d")

FEED.addEventListener("click", () => {
    currObjective = "feed"
    goTo = [randomCoords(), randomCoords()];
});

PLAY.addEventListener("click", () => {
    currObjective = "play"
});

SLEEP.addEventListener("click", () => {
    currObjective = "sleep"
    goTo = [randomCoords(), randomCoords()]
})

setInterval(() => {
    if (connected && sprite && ctx) {
        ctx.clearRect(0, 0, CANVAS.width, CANVAS.height);
        const targetX = goTo[0];
        const targetY = goTo[1];

        const dx = targetX - x;
        const dy = targetY - y;

        const distance = Math.sqrt(dx * dx + dy * dy);
        let rotate = false

        if (distance > 8) {

            const stepSize = 0.01;
            x += (dx / distance) * stepSize;
            y += (dy / distance) * stepSize;

            const wobbleRange = 0.02;
            x += (Math.random() - 0.5) * wobbleRange;
            y += (Math.random() - 0.5) * wobbleRange;


            const jumpProbability = 0.02;
            if (Math.random() < jumpProbability) {
                const jumpSize = 2.0;
                x += (Math.random() - 0.5) * jumpSize;
                y += (Math.random() - 0.5) * jumpSize;
            }
        } else {
            switch (currObjective) {
                case "wander":
                    goTo = [randomCoords(), randomCoords()];
                    break;
                case "play":
                    const jumpProbability = 0.05;
                    if (Math.random() < jumpProbability) {
                        const jumpHeight = 5.0;
                        x += (Math.random() - 0.5) * jumpHeight;
                        y += (Math.random() - 0.5) * jumpHeight;
                        break;
                    }
                    let playResponse: RequestObject = {
                        type: "Feed",
                        data: f32ToBytes(0.05)
                    };
                    sendQuery(playResponse);
                    break;
                case "feed":
                    let foodResponse: RequestObject = {
                        type: "Feed",
                        data: f32ToBytes(10.0)
                    };
                    sendQuery(foodResponse);
                    currObjective = "wander"
                    goTo = [randomCoords(), randomCoords()];
                    break;
                case "sleep":
                    rotate = true
                    let sleepyResponse: RequestObject = {
                        type: "Sleep",
                        data: f32ToBytes(0.5)
                    };
                    sendQuery(sleepyResponse);
            }
        }

        const imageData = ctx.createImageData(16, 16);

        for (let i = 0; i < sprite.length; i++) {
            imageData.data[i] = sprite[i];
        }

        if (rotate) {
            let rotated = rotateImageData90Degrees(imageData.data, 16, 16)
            for (let i = 0; i < rotated.length; i++) {
                imageData.data[i] = rotated[i];
            }
        }


        // Draw the sprite at the updated position
        ctx.putImageData(imageData, x, y);
        switch (currObjective) {
            case "play":
                drawRedBall(ctx, goTo[0], goTo[1], 3);
                break;
            case "feed":
                ctx.fillRect(goTo[0], goTo[1], 2, 2);
                break;
            case "sleep":
                ctx.fillRect(goTo[0], goTo[1], 7, 4);
                ctx.fillStyle = "white";
                ctx.fillRect(goTo[0] + 1, goTo[1] + 1, 5, 2);
                ctx.fillStyle = "black"
                break;
        }


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

function randomCoords() {
    return Math.random() * 32
}

function rotateImageData90Degrees(sprite: Uint8ClampedArray, width: number, height: number): Uint8ClampedArray {
    const rotatedData = new Uint8ClampedArray(sprite.length);

    for (let y = 0; y < height; y++) {
        for (let x = 0; x < width; x++) {
            const srcIndex = (y * width + x) * 4;

            const newX = height - 1 - y;
            const newY = x;
            const destIndex = (newY * height + newX) * 4;

            rotatedData[destIndex] = sprite[srcIndex];       // R
            rotatedData[destIndex + 1] = sprite[srcIndex + 1]; // G
            rotatedData[destIndex + 2] = sprite[srcIndex + 2]; // B
            rotatedData[destIndex + 3] = sprite[srcIndex + 3]; // A
        }
    }

    return rotatedData;
}

function drawRedBall(ctx: CanvasRenderingContext2D, x: number, y: number, radius: number) {
    ctx.fillStyle = "red";
    ctx.beginPath();
    ctx.arc(x, y, radius, 0, Math.PI * 2);
    ctx.fill();
    ctx.closePath();
    ctx.fillStyle = "black";
}