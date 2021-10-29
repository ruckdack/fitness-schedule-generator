import { animate, style, transition, trigger } from '@angular/animations';
import { Component, Input, OnInit } from '@angular/core';
import { Subscriber } from 'rxjs';
import { ConfigService } from '../services/config/config.service';
import { Exercise } from '../services/config/types.config.service';

@Component({
  selector: 'app-exercise',
  templateUrl: './exercise.component.html',
  styleUrls: ['./exercise.component.sass'],
  animations: [
    trigger('foldInOut', [
      transition(':enter', [
        style({ height: 0, 'padding-top': 0, 'margin-bottom': 0 }),
        animate(100),
      ]),
      transition(':leave', [
        animate(
          100,
          style({ height: 0, 'padding-top': 0, 'margin-bottom': 0 })
        ),
      ]),
    ]),
  ],
})
export class ExerciseComponent implements OnInit {
  @Input() exercise!: Exercise;
  @Input() pushExercises!: () => void;
  @Input() expand!: () => void;
  @Input() close!: () => void;
  @Input() expanded!: boolean;

  public muscleNames: string[];

  constructor(public configService: ConfigService) {
    this.expanded = false;
    this.muscleNames = [];
  }

  ngOnInit() {
    this.configService.subMuscleNames(
      Subscriber.create((muscleNames?: string[]) => {
        this.muscleNames = muscleNames!;
      })
    );
  }

  toggleExpansion() {
    if (this.expanded) {
      this.close();
      this.expanded = false;
      return;
    }
    this.expand();
    this.expanded = true;
  }

  update1RM(new1RM: any) {
    this.exercise.initial_1rm = new1RM;
    this.pushExercises();
  }

  updateReps(newReps: any) {
    this.exercise.reps = newReps;
    this.pushExercises();
  }

  updateTarget(newTarget: any) {
    console.log(newTarget);
    this.exercise.target = newTarget;
    this.pushExercises();
  }
}
