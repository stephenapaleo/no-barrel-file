import { CircularB } from "./circular-b";

export class CircularA {
  dep = new CircularB();
}
