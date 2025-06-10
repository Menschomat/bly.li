import { Component, EventEmitter, Output } from '@angular/core';

@Component({
  selector: 'app-short-table-row',
  imports: [],
  host: {
    class: 'flex-1 grid grid-cols-[6rem_1fr_5rem_auto] grid-rows-1 gap-4',
  },
  template: `
    <div>
      <i
        class="fa-regular fa-copy cursor-pointer"
        (click)="handleCopyClick($event)"
      ></i
      >&nbsp;&nbsp;<ng-content select="row-title"></ng-content>
    </div>
    <div class="truncate overflow-hidden text-ellipsis whitespace-nowrap">
      <ng-content select="row-url"></ng-content>
    </div>
    <div
      class="flex items-center max-w-fit bg-purple-100 text-purple-800 text-xs font-medium me-2 px-2.5 py-0.5 rounded-sm dark:bg-purple-900 dark:text-purple-300"
    >
      <ng-content select="row-count"></ng-content>
    </div>
    <div>
      <i
        class="fa-regular fa-trash-can cursor-pointer"
        (click)="handleDeleteClick($event)"
      ></i>
    </div>
  `,
  styles: ``,
})
export class ShortTableRowComponent {
  @Output()
  public delete = new EventEmitter<MouseEvent>();
  @Output()
  public copy = new EventEmitter<MouseEvent>();

  public handleDeleteClick(event: MouseEvent) {
    this.delete.emit(event);
  }

  public handleCopyClick(event: MouseEvent) {
    this.copy.emit(event);
  }
}
