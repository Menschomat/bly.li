import { Component } from '@angular/core';

@Component({
  selector: 'app-menu-item',
  imports: [],
  template: `
    <li>
      <a
        class="block px-4 py-2 bg-transparent dark:hover:text-white cursor-pointer  hover:bg-gray-200/50 hover:dark:bg-gray-800/50"
        ><ng-content select="item-content"></ng-content
      ></a>
    </li>
  `,
  styles: ``,
})
export class MenuItemComponent {}
