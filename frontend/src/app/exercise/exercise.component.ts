import {
  animate,
  state,
  style,
  transition,
  trigger,
} from '@angular/animations';
import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-exercise',
  templateUrl: './exercise.component.html',
  styleUrls: ['./exercise.component.sass'],
  animations: [
    trigger('flyInOut', [
      transition(':enter', [
        style({ height: 0, color: 'transparent' }),
        animate(100),
      ]),
      transition(':leave', [
        animate(100, style({ height: 0, color: 'transparent' })),
      ]),
    ]),
  ],
})
export class ExerciseComponent implements OnInit {
  @Input() exercise!: any;
  public expanded: boolean;

  constructor() {
    this.expanded = false;
  }

  ngOnInit() {}

  toggleExpansion() {
    this.expanded = !this.expanded;
  }
}
