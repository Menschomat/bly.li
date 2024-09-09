import { BehaviorSubject } from 'rxjs';
import { Injectable } from '@angular/core';
import { ShortURL } from '../core/api/v1';

@Injectable({
  providedIn: 'root',
})
export class URLService {
  private lastShortUrl = new BehaviorSubject<string | undefined>('');
  private ownUrls = new BehaviorSubject<Array<ShortURL> | undefined>([]);
  constructor() {}
  public triggerNextOwnUrls(url: Array<ShortURL> | undefined) {
    this.ownUrls.next(url);
  }
  public triggerNextShort(url: string | undefined) {
    this.lastShortUrl.next(url);
  }
  public clearLastShort() {
    this.lastShortUrl.next(undefined);
  }
  get lastShortUrl$() {
    return this.lastShortUrl.asObservable();
  }
  get ownUrls$() {
    return this.ownUrls.asObservable();
  }
}
