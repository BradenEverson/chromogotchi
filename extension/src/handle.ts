type PetId = string;
type ErrorString = string;

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
});

socket.addEventListener("message", (event) => {
  console.log(event.data);
  const data = JSON.parse(event.data);
  console.log(data);
});

socket.addEventListener("close", () => {
  console.log("Disconnected from the WebSocket server.");
});

socket.addEventListener("error", (error) => {
  console.error("WebSocket error:", error);
});

window.addEventListener("beforeunload", () => {
  socket.close();
});

let id = getCookie("PETID");
if (!id) {
  id = await newPet();
}

console.log(id);

let hunger: number, happiness: number, sleepyness: number;
let sprite: number[];
let lilGuyX: number,
  lilGuyY: number = 0;

setInterval(() => {
  /// Get Hunger
  /// Get Happiness
  /// Get Sleepyness
  /// Get Necessary Sprite
}, 1000);

setInterval(() => {
  /// Update the dude's movement
}, 250);
