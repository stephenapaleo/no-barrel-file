import { CircularB } from ".";

export class CircularA {
  dep = new CircularB();
}
