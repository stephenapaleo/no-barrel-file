import { CircularA } from ".";

export class CircularB {
  dep = new CircularA();
}
