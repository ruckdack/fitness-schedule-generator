import { Component } from '@angular/core';
import { AuthService } from './services/auth/auth.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.sass'],
})
export class AppComponent {
  public show = false;

  constructor(authService: AuthService) {
    authService.userState.subscribe((isLoggedIn) => {
      this.show = isLoggedIn;
    });
  }
}
