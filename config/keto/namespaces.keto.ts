import { Context, Namespace } from "@ory/keto-namespace-types"

class User implements Namespace {}

class Role implements Namespace {
  related: {
    members: User[]
  }
}

/**
 * "Product" is a namespace representing a product. It has some rewrites.
 */
class Product implements Namespace {
  // Relations are defined and type-annotated.
  related: {
    /**
     * "owners" are the users that are the owners of the product.
     */
    owners: User[]
    /**
     * "admins" are the roles that are administrators of this product (potentially only one).
     */
    admins: Role[]
    /**
     * "special_roles" are the roles a user has to be member of to gain "additional_permissions"
     */
    special_roles: Role[]
  }

  permits = {
    // this is probably three/four rewrites (create, read, update, delete) with similar rules
    crud: (ctx: Context): boolean =>
      this.related.owners.includes(ctx.subject) ||
      this.related.admins.traverse((admin) => admin.related.members.includes(ctx.subject)),

    // for the additional_permissions one has to have crud and be member of a special role
    additional_permissions: (ctx: Context): boolean =>
      this.permits.crud(ctx) &&
      this.related.special_roles.traverse((role) => role.related.members.includes(ctx.subject))
  }
}

class Osu implements Namespace {
    related: {
        /**
         * Tournament roles
         */
        host: Role[]
        
        // General roles
        designer: Role[]
        
        // Mappool roles
        map_pooler: Role[]
        mapper: Role[]
        testplayer: Role[]

        // Match roles
        Refree: Role[]
        Stremer: Role[]
        Commentator: Role[]
    }
    
    permits = {
        
    }
}
