import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { MatIcon } from '@angular/material/icon';
import { ActivatedRoute, Router } from '@angular/router';
import { CollectionSearchResponse, CollectionSearchService } from '../../services/collection-search.service';
import { SearchResultItemComponent } from "../search-result-item/search-result-item.component";

import {MatFormFieldModule} from '@angular/material/form-field';
import {MatIconModule} from '@angular/material/icon';
import {MatInputModule} from '@angular/material/input';

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [
    CommonModule,
    MatIcon,
    ReactiveFormsModule,
    SearchResultItemComponent,
    MatFormFieldModule,
    MatIcon,
    MatInputModule,
  ],
  templateUrl: './search.component.html',
  styleUrl: './search.component.css'
})
export class SearchComponent {

  @Input()collection: string = "no collection"
  searchQuery = new FormControl("")
  searchResult: CollectionSearchResponse | undefined

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private searchService: CollectionSearchService
  ) { }
  


  search() {
    let query = this.searchQuery.value!
    console.log("Query: " + query);
    this.searchService.searchCollection(this.collection,query).subscribe(data => {
      console.log(data);
      
      this.searchResult = data
    }
    );
  }

}
