import { Component, OnInit } from '@angular/core';
import { Observable, Subscription } from 'rxjs';
import { AuthService } from '../services/auth/auth.service';

@Component({
  selector: 'app-auth',
  templateUrl: './auth.component.html',
  styleUrls: ['./auth.component.sass'],
})
export class AuthComponent implements OnInit {
  showLoginButton: boolean;
  showLogoutButton: boolean;
  userStateSub: Subscription;
  userID: string;

  constructor(private authService: AuthService) {
    this.showLoginButton = false;
    this.showLogoutButton = false;
    this.userStateSub = Subscription.EMPTY;
    this.userID = '';
  }

  ngOnInit(): void {
    this.authService.userState.subscribe((isLoggedIn) => {
      this.showLoginButton = !isLoggedIn;
      this.showLogoutButton = isLoggedIn;
      this.userID = this.authService.userID;
    });
  }

  ngOnDestroy(): void {
    if (this.userStateSub) {
      this.userStateSub.unsubscribe();
    }
  }

  login(): void {
    this.authService.login();
  }

  logout(): void {
    this.authService.logout();
  }
}
