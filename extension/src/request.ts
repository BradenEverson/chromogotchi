export type RequestObject = {
  type: string;
  data: number[];
};

export type ResponseObject = {
  type: string;
  // Base64 encoded
  data: string;
};

export type Objective = [number, number] | "sleep" | "eat" | ["play", [number, number]]