import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { ActivatedRoute, ParamMap, Router } from '@angular/router';
import { switchMap } from 'rxjs';

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './search.component.html',
  styleUrl: './search.component.css'
})
export class SearchComponent {

  collectionName: string = "no collection"

  constructor(
    private route: ActivatedRoute,
    private router: Router,
  ) { }
  
  @Input()
  set collection(collection: string) {
    this.collectionName = collection
  }

  
//   ngOnInit() {
//    this.route.paramMap.pipe(
//     switchMap((params: ParamMap) =>
//       this.collection = params.get('collection')!
//   ));
//  }

}
