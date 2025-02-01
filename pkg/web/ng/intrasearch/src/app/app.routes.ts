import { Routes } from '@angular/router';
import { CollectionsListComponent } from './components/collections-list/collections-list.component';
import { SearchComponent } from './components/search/search.component';

export const routes: Routes = [
  { path: '', component: CollectionsListComponent },
    { path: 'search', component: SearchComponent },
];
