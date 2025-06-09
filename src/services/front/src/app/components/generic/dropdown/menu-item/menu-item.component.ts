import { Component } from '@angular/core';

@Component({
  selector: 'app-menu-item',
  imports: [],
  template: `
    <li>
      <a
        class="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white cursor-pointer"
        ><ng-content select="item-content"></ng-content
      ></a>
    </li>
  `,
  styles: ``,
})
export class MenuItemComponent {}
