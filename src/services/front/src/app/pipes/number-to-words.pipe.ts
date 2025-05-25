import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'numberToWords',
})
export class NumberToWordsPipe implements PipeTransform {
  transform(value: number): string {
    if (value === 0) {
      return 'zero';
    }
    return this.convertNumberToWords(value);
  }

  private convertNumberToWords(num: number): string {
    const units = [
      '',
      'one',
      'two',
      'three',
      'four',
      'five',
      'six',
      'seven',
      'eight',
      'nine',
    ];
    const teens = [
      'ten',
      'eleven',
      'twelve',
      'thirteen',
      'fourteen',
      'fifteen',
      'sixteen',
      'seventeen',
      'eighteen',
      'nineteen',
    ];
    const tens = [
      '',
      'ten',
      'twenty',
      'thirty',
      'forty',
      'fifty',
      'sixty',
      'seventy',
      'eighty',
      'ninety',
    ];

    if (num < 10) return units[num];
    if (num >= 10 && num < 20) return teens[num - 10];
    if (num >= 20 && num < 100)
      return (
        tens[Math.floor(num / 10)] + ' ' + this.convertNumberToWords(num % 10)
      );
    if (num >= 100 && num < 1000)
      return (
        units[Math.floor(num / 100)] +
        ' hundred ' +
        this.convertNumberToWords(num % 100)
      );
    if (num >= 1000 && num < 1000000)
      return (
        this.convertNumberToWords(Math.floor(num / 1000)) +
        ' thousand ' +
        this.convertNumberToWords(num % 1000)
      );
    if (num >= 1000000 && num < 1000000000)
      return (
        this.convertNumberToWords(Math.floor(num / 1000000)) +
        ' million ' +
        this.convertNumberToWords(num % 1000000)
      );
    return 'Number is too large';
  }
}
