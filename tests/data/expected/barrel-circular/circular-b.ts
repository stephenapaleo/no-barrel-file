import { CircularA } from "./circular-a";

export class CircularB {
  dep = new CircularA();
}
