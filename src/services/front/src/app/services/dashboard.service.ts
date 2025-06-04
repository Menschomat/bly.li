import { Injectable } from '@angular/core';
import { DasherService, ShortClickCount, ShortURL } from '../api';
import { BehaviorSubject, Observable, Subject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class DashboardService {
  private readonly allShorts: Subject<ShortURL[]> = new BehaviorSubject(
    [] as ShortURL[]
  );
  constructor(private readonly dasherService: DasherService) {
    this.refresh();
  }

  public get $allShorts(): Observable<ShortURL[]> {
    return this.allShorts.asObservable();
  }
  public delete(short: ShortURL) {
    this.dasherService
      .dasherShortShortDelete(short.Short ?? '')
      .subscribe(() => this.refresh());
  }

  public clickHistoryByShortUrl(
    short: ShortURL
  ): Observable<ShortClickCount[]> {
    return this.clickHistoryByShort(short.Short ?? '');
  }

  public clickHistoryByShort(short: string): Observable<ShortClickCount[]> {
    console.log("SHORT",short);
    
    return this.dasherService.dasherShortShortClicksGet(short ?? '');
  }

  public refresh() {
    this.dasherService
      .dasherShortAllGet()
      .subscribe((shorts) => this.allShorts.next(shorts));
  }
}
