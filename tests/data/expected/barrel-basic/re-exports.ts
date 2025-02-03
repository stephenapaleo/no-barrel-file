import { BasicClass } from "./classes";

export { BasicClass as RenamedBasicClass };

import { BasicType } from "./types";

type RexportedBasicType = BasicType;
export type { RexportedBasicType };

const RexportedBasicConstToExport = 42;
export { RexportedBasicConstToExport };
