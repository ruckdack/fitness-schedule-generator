import { Injectable } from '@angular/core';
import {
  Auth,
  authState,
  signOut,
  User,
  GoogleAuthProvider,
  signInWithPopup,
} from '@angular/fire/auth';
import { EMPTY, Observable } from 'rxjs';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  user: Observable<User | null> = EMPTY;
  userState: Observable<boolean> = EMPTY;
  userID: string = '';

  constructor(private auth: Auth) {
    if (auth) {
      this.user = authState(this.auth);
      this.userState = this.user.pipe(map((u) => !!u));
      this.user.subscribe((user) => {
        this.userID = user ? user.uid : '';
      });
    }
  }

  async login() {
    return await signInWithPopup(this.auth, new GoogleAuthProvider());
  }

  async logout() {
    return await signOut(this.auth);
  }
}
