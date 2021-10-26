import { Component, OnInit } from '@angular/core';
import { Subscriber, Subscription } from 'rxjs';
import { ConfigService } from '../services/config/config.service';
import { Config, Exercise } from '../services/config/types.config.service';

@Component({
  selector: 'app-exercise-list',
  templateUrl: './exercise-list.component.html',
  styleUrls: ['./exercise-list.component.sass'],
})
export class ExerciseListComponent implements OnInit {
  public exercises: Exercise[];
  private exerciseSubscription: Subscription = Subscription.EMPTY;

  constructor(public configService: ConfigService) {
    this.exercises = [];
  }

  ngOnInit(): void {
    this.exerciseSubscription = this.configService.subExercises(
      Subscriber.create((exercises?: Exercise[]) => {
        this.exercises = exercises!;
      })
    );
  }

  ngOnDestroy(): void {
    this.exerciseSubscription.unsubscribe();
  }
}
