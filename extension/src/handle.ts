type PetId = string;
type ErrorString = string;

const HTTP_LOCATION: string = "http://localhost:7878"

export function getCookie(name: string): PetId | null {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop()?.split(';').shift() || null;
    return null;
}

export async function newPet(): Promise<PetId | ErrorString> {
    try {
        const response = await fetch(HTTP_LOCATION + "/newpet", {
            method: 'GET'
        });

        if (response.status === 200) {
            const data = await response.json();
            const petID = data.PetId;
            document.cookie = `PETID=${petID}; path=/;`;
            return petID as PetId;
        } else {
            return "Bad response from endpoint :(" as ErrorString
        }
    } catch (err) {
        return "Error reaching HTTP endpoint" as ErrorString
    }
}
