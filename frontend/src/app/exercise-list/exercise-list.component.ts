import { Component, OnInit } from '@angular/core';
import { Subscriber, Subscription } from 'rxjs';
import { ConfigService } from '../services/config/config.service';
import { Exercise } from '../services/config/types.config.service';

@Component({
  selector: 'app-exercise-list',
  templateUrl: './exercise-list.component.html',
  styleUrls: ['./exercise-list.component.sass'],
})
export class ExerciseListComponent implements OnInit {
  public exercises: Exercise[];
  public pushExercises: () => void;
  public currentlyExpanded: number | null;
  private exercisesSubscription: Subscription;

  constructor(public configService: ConfigService) {
    this.exercises = [];
    this.currentlyExpanded = null;
    this.pushExercises = () => {};
    this.exercisesSubscription = Subscription.EMPTY;
  }

  ngOnInit(): void {
    this.exercisesSubscription = this.configService.subExercises(
      Subscriber.create((exercises?: Exercise[]) => {
        if (this.exercises.length !== exercises!.length) {
          this.exercises = exercises!;
          return;
        }
        this.exercises.forEach((exercise: Exercise, idx: number) =>
          Object.assign(exercise, exercises![idx])
        );
      })
    );

    this.pushExercises = () => this.configService.pushExercises(this.exercises);
  }

  ngOnDestroy(): void {
    this.exercisesSubscription.unsubscribe();
  }

  expandWrapper(idx: number) {
    return () => (this.currentlyExpanded = idx);
  }

  closeWrapper() {
    return () => (this.currentlyExpanded = null);
  }
}
