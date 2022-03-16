import * as Knex from "knex";


export async function up(knex: Knex): Promise<void> {
    return knex.schema.table("vcs", function (table) {
        table.dropColumn("credential_type")
    })
}

export async function down(knex: Knex): Promise<void> {
    return knex.schema.table("vcs", function (table) {
        table.string('credential_type', 255).notNullable()
    })
}

