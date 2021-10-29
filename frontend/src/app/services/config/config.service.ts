import { Injectable } from '@angular/core';
import {
  Firestore,
  onSnapshot,
  DocumentReference,
  doc,
  setDoc,
  updateDoc,
  query,
} from '@angular/fire/firestore';
import { EMPTY, Observable, Subscriber } from 'rxjs';
import { map } from 'rxjs/operators';
import { AuthService } from '../auth/auth.service';
import { Config, Exercise } from './types.config.service';

@Injectable({
  providedIn: 'root',
})
export class ConfigService {
  private initalConfig: Config = JSON.parse(
    `{    "weekdays": {        "mon": "upper",        "tue": "lower",        "wed": "upper",        "fri": "lower",        "sat": "upper"    },    "start_date": "2021-10-16",    "splits": [        {            "name": "upper",            "exercises": [                [                    ["incline benchpress (dumbbell)", "benchpress"],                    ["high row (pulley)"]                ],                [                    ["lateral raise (dumbbell)", "lateral raise (machine)"],                    ["lat pulldown"]                ],                [                    ["incline curl (dumbbell)", "curl machine"],                    ["triceps extension", "overhead triceps extension"]                ]            ]        },        {            "name": "lower",            "exercises": [                [["leg press"]],                [["RDL", "hip thrust"]],                [["leg extension"], ["laying leg curl"]],                [["seated calf raise"], ["hyperextensions"]],                [["standing calf raise"], ["crunch machine"]]            ]        }    ],    "exercises": [        {            "name": "incline benchpress (dumbbell)",            "initial_1rm": 10,            "reps": 9,            "target": "chest"        },        {            "name": "benchpress",            "initial_1rm": 200,            "reps": 8,            "target": "chest"        },        {            "name": "high row (pulley)",            "initial_1rm": 300,            "reps": 8,            "target": "lower traps / rhomboids"        },        {            "name": "lateral raise (dumbbell)",            "initial_1rm": 400,            "reps": 8,            "target": "delts"        },        {            "name": "lateral raise (machine)",            "initial_1rm": 500,            "reps": 8,            "target": "delts"        },        {            "name": "lat pulldown",            "initial_1rm": 600,            "reps": 8,            "target": "lat"        },        {            "name": "incline curl (dumbbell)",            "initial_1rm": 700,            "reps": 8,            "target": "biceps"        },        {            "name": "curl machine",            "initial_1rm": 800,            "reps": 8,            "target": "biceps"        },        {            "name": "triceps extension",            "initial_1rm": 900,            "reps": 8,            "target": "triceps"        },        {            "name": "overhead triceps extension",            "initial_1rm": 1000,            "reps": 8,            "target": "triceps"        },        {            "name": "crunch machine",            "initial_1rm": 1100,            "reps": 8,            "target": "abs"        },        {            "name": "standing calf raise",            "initial_1rm": 1200,            "reps": 8,            "target": "calves"        },        {            "name": "seated calf raise",            "initial_1rm": 1300,            "reps": 8,            "target": "calves"        },        {            "name": "laying leg curl",            "initial_1rm": 1400,            "reps": 8,            "target": "hamstrings"        },        {            "name": "hip thrust",            "initial_1rm": 1500,            "reps": 8,            "target": "glute"        },        {            "name": "RDL",            "initial_1rm": 1600,            "reps": 8,            "target": "hamstrings"        },        {            "name": "leg extension",            "initial_1rm": 1700,            "reps": 8,            "target": "quads"        },        {            "name": "leg press",            "initial_1rm": 1800,            "reps": 8,            "target": "quads"        },        {            "name": "hyperextensions",            "initial_1rm": 50,            "reps": 12,            "target": "lower back"        }    ],    "muscles": [        {            "name": "quads",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "chest",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "hamstrings",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "delts",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "calves",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "glute",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "lower traps / rhomboids",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "lat",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "delts",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "abs",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "biceps",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "triceps",            "sets": [12, 14, 14, 16, 8]        },        {            "name": "lower back",            "sets": [12, 14, 14, 16, 8]        }    ]}`
  );
  private ref: DocumentReference;
  private observableConfig: Observable<any> = EMPTY;

  constructor(firestore: Firestore, authService: AuthService) {
    this.ref = doc(firestore, 'config', authService.userID);
    this.observableConfig = new Observable((sub: Subscriber<Config>) => {
      onSnapshot(this.ref, (snapshot) => {
        if (snapshot.exists()) {
          const config: Config = this.rebuildArrays(snapshot.data());
          sub.next(config);
          return;
        }
        this.setInitialConfig();
      });
    });
  }

  setInitialConfig() {
    setDoc(this.ref, this.removeArrays(this.initalConfig));
  }

  subExercises(sub: Subscriber<Exercise[]>) {
    return this.observableConfig
      .pipe(map((config: Config) => config.exercises))
      .subscribe(sub);
  }

  subMuscleNames(sub: Subscriber<string[]>) {
    return this.observableConfig
      .pipe(
        map((config: Config) => config.muscles.map((muscle) => muscle.name))
      )
      .subscribe(sub);
  }

  pushExercises(exercises: Exercise[]) {
    updateDoc(this.ref, { exercises: this.removeArrays(exercises) });
  }

  removeArrays(obj: any) {
    // checks if primitive
    if (obj !== Object(obj)) {
      return obj;
    }
    let newObj: any = {};
    if (Array.isArray(obj)) {
      obj.forEach((arrElement: any, idx: number) => {
        newObj[idx.toString(10)] = this.removeArrays(arrElement);
      });
      return newObj;
    }
    Object.entries(obj).forEach(([k, v]) => {
      newObj[k] = this.removeArrays(v);
    });
    return newObj;
  }

  rebuildArrays(obj: any) {
    const isArray = () => {
      // if key not parseable => definetly not an array
      if (Object.keys(obj).some((key) => isNaN(Number(key)))) return false;
      const indicees = Object.keys(obj)
        .map((i) => parseInt(i))
        .sort((a, b) => a - b);
      return indicees.every((val, idx) => val === idx);
    };
    // checks if primitive
    if (obj !== Object(obj)) {
      return obj;
    }
    if (isArray()) {
      let arr: any[] = [];
      Object.entries(obj).forEach(([k, v]) => {
        arr[parseInt(k)] = this.rebuildArrays(v);
      });
      return arr;
    }
    let newObj: any = {};
    Object.entries(obj).forEach(([k, v]) => {
      newObj[k] = this.rebuildArrays(v);
    });
    return newObj;
  }
}
