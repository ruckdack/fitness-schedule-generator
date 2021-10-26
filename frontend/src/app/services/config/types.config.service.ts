export type Exercise = {
  initial_1rm: number;
  name: string;
  reps: number;
  target: string;
};

export type Muscle = {
  name: string;
  sets: number[];
};

export type Variation = string[];

export type Superset = Variation[];

export type Split = {
  name: string;
  exercises: Superset[];
};

export type Config = {
  weekdays: { [key: string]: string };
  start_date: string;
  split: Split[];
  exercises: Exercise[];
  muscles: Muscle[];
};
