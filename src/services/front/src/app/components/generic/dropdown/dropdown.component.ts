import { CommonModule } from '@angular/common';
import { Component, ElementRef, HostListener } from '@angular/core';

@Component({
  selector: 'app-dropdown',
  imports: [CommonModule],
  template: `
    <button
      (click)="open = !open"
      class="
        flex items-center 
        cursor-pointer 
        rounded-lg 
        pe-1 
        text-lg 
        text-gray-900 
        hover:text-indigo-600 
        md:me-0 
        focus:outline-hidden 
        focus:ring-0 
        focus:ring-gray-100 
        dark:text-white 
        dark:hover:text-indigo-500 
        dark:focus:ring-gray-700
      "
      type="button"
    >
      <span class="sr-only">Open user menu</span>
      <i class="fa-regular fa-user mr-3"></i>
      <ng-content select="user-name"></ng-content>
      <i class="ml-2 fa-solid fa-chevron-down"></i>
    </button>
    @if (open == true) {
    <!-- Dropdown menu -->
    <div
      id="dropdownAvatarName"
      class="
        z-10 
        absolute 
        mt-2 
        w-38 
        bg-white 
        rounded-2xl 
        shadow 
        shadow-sm 
        border 
        border-gray-100 
        dark:bg-black 
        dark:border-gray-900
      "
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
  constructor(private readonly elementRef: ElementRef) {}

  @HostListener('document:click', ['$event'])
  onDocumentClick(event: MouseEvent) {
    const clickedInside = this.elementRef.nativeElement.contains(event.target);
    if (!clickedInside && this.open) {
      this.open = false;
    }
  }
}
