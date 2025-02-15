import { Injectable } from '@angular/core';
import { DasherService, ShortURL } from '../api';
import { BehaviorSubject, Observable, Subject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class DashboardService {
  private allShorts: Subject<ShortURL[]> = new BehaviorSubject(
    [] as ShortURL[]
  );
  constructor(private dasherService: DasherService) {
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

  public refresh() {
    this.dasherService
      .dasherShortAllGet()
      .subscribe((shorts) => this.allShorts.next(shorts));
  }
}
