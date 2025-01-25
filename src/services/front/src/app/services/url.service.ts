import { BehaviorSubject } from 'rxjs';
import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class URLService {
  private lastShortUrl = new BehaviorSubject<string | undefined>('');
  constructor() {}
  public triggerNextShort(url: string | undefined) {
    this.lastShortUrl.next(url);
  }
  public clearLastShort() {
    this.lastShortUrl.next(undefined);
  }
  get lastShortUrl$() {
    return this.lastShortUrl.asObservable();
  }
}
