export type RequestObject = {
  type: string;
  data: number[];
};

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
