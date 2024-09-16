export type RequestObject = {
  type: string;
  data: number[];
};

export type ResponseObject = {
  type: string;
  // Base64 encoded
  data: string;
};

export type Objective = "sleep" | "play" | "feed" | "wander"