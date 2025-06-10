import { CommonModule } from '@angular/common';
import { Component, ElementRef, HostListener } from '@angular/core';

@Component({
  selector: 'app-dropdown',
  imports: [CommonModule],
  template: `
    <button
      (click)="open = !open"
      class="flex items-center text-lg pe-1  text-gray-900 focus:outline-hidden rounded-lg cursor-pointer hover:text-indigo-600 dark:hover:text-indigo-500 md:me-0 focus:ring-0 focus:ring-gray-100 dark:focus:ring-gray-700 dark:text-white"
      type="button"
    >
      <span class="sr-only">Open user menu</span>
      <i class=" fa-regular fa-user mr-3"></i>
      <ng-content select="user-name"></ng-content>
      <i class="ml-2 fa-solid fa-chevron-down"></i>
    </button>
    @if (open == true) {
    <!-- Dropdown menu -->
    <div
      id="dropdownAvatarName"
      class="z-10 absolute mt-2 divide-y divide-gray-100 rounded-lg shadow-sm w-44 dark:divide-gray-600 backdrop-blur-md bg-white/30 shadow dark:bg-gray-800/50 rounded-2xl backdrop-blur-md"
    >
      <ul class="py-2 text-sm text-gray-700 dark:text-gray-200">
        <ng-content select="item-list"></ng-content>
      </ul>
    </div>
    }
  `,
  styles: ``,
})
export class DropdownComponent {
  public open: boolean = false;
  constructor(private elementRef: ElementRef) {}

  @HostListener('document:click', ['$event'])
  onDocumentClick(event: MouseEvent) {
    const clickedInside = this.elementRef.nativeElement.contains(event.target);
    if (!clickedInside && this.open) {
      this.open = false;
    }
  }
}
